package com.example.gdsudao

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.widget.Toast
import kotlinx.android.synthetic.main.activity_detalhes.*

class DetalhesActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_detalhes)

        setSupportActionBar(findViewById(R.id.toolbarDetalhes))

        // Get the areas
        var sp = com.example.gdsudao.utils.SharedPreferences()
        var areas = sp.RecuperarListaAreas(this)

        val bundle = intent.extras
        val areaIndex = bundle?.getInt("item")

        if(areaIndex != null){
            // Toast.makeText(this, "Area: ${areaIndex}", Toast.LENGTH_SHORT).show(); FUNCIONANDO
            Log.println(Log.DEBUG, "area", "Area: ${areas[areaIndex].toString()}")

            tvDataEstimada.text = areas[areaIndex].proxcorte
            tvNomeArea.text = areas[areaIndex].nome
            tvDataUltimoPastejo.text = areas[areaIndex].dataCorte
            tvGdAcumulado.text = areas[areaIndex].st
            tvEstacao.text  = areas[areaIndex].codigoEstacao
            tvPrevisoes.text  = areas[areaIndex].previsao + "%"
            tvNormais.text  = areas[areaIndex].normal + "%"
            tvDiarios.text = areas[areaIndex].diario + "%"
            tvNumeroCortes.text = areas[areaIndex].numeroCorte

        }else{
            Toast.makeText(this, "Erro ao apresentar os detalhes.", Toast.LENGTH_SHORT).show()
        }

    }
}